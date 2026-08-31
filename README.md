# Base configuration role for Ansible

[![molecule](https://github.com/megalomania428/ansible-mega-base/actions/workflows/role-test.yaml/badge.svg)](https://github.com/megalomania428/ansible-mega-base/actions/workflows/role-test.yaml)

The role performs a base host configuration

## Role release to Ansible galaxy

- clone me:

```bash
git clone --recursive git@github.com:megalomania428/ansible-mega-base.git ansible-mega-base
```

- make tag and send to release:

```bash
git checkout master && git pull
git tag -fm $(git branch --sho) 1.0.0 && git push --force origin $(git describe)
```
