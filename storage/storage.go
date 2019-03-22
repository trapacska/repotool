package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const file = "./.storage.json"

type Table struct {
	GithubAccessToken       string
	TransferredRepositories []Repo
	UpdatedRepositories     []string
}

type Repo struct {
	OriginalURL string
	TargetOrg   string
}

func open() (Table, error) {
	var t Table

	data, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
		return t, err
	}

	return t, json.Unmarshal(data, &t)
}

func (t Table) save() error {
	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0644)
}

func Read(fn func(t *Table)) error {
	table, err := open()
	if err != nil {
		return err
	}

	fn(&table)

	return nil
}

func Update(fn func(t *Table)) error {
	table, err := open()
	if err != nil {
		return err
	}

	fn(&table)

	return table.save()
}
