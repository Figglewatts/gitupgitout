package config

import (
	"errors"
	"fmt"
	"io"

	"github.com/figglewatts/gitupgitout/pkg/source"
	"github.com/figglewatts/gitupgitout/pkg/source/github"
	"github.com/figglewatts/gitupgitout/pkg/source/gitlab"
	"gopkg.in/yaml.v3"
)

type Mirror struct {
	Source struct {
		GitHubAccount *github.Account `yaml:"githubAccount"`
		GitLabAccount *gitlab.Account `yaml:"gitlabAccount"`
	} `yaml:"source"`
	CloneTo string `yaml:"cloneTo"`
}

func (m Mirror) GetSource() source.RepoSource {
	if m.Source.GitHubAccount != nil {
		return m.Source.GitHubAccount
	} else if m.Source.GitLabAccount != nil {
		return m.Source.GitLabAccount
	}
	return nil
}

func (m Mirror) validate() error {
	if firstNonNil(
		m.Source.GitHubAccount,
		m.Source.GitLabAccount,
	) == nil {
		return errors.New("no source")
	}

	return nil
}

type Config struct {
	Mirrors []Mirror `yaml:"mirrors"`
}

func (c Config) validate() error {
	if len(c.Mirrors) == 0 {
		return errors.New("mirrors must have at least one entry")
	}

	for i, m := range c.Mirrors {
		if err := m.validate(); err != nil {
			return fmt.Errorf("mirror %d: %w", i, err)
		}
	}

	return nil
}

func Load(configReader io.Reader) (*Config, error) {
	var conf Config
	decoder := yaml.NewDecoder(configReader)
	err := decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}

	if err = conf.validate(); err != nil {
		return nil, err
	}

	return &conf, nil
}
