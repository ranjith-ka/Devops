### Release Process for Project Delivery

Upstream --> Branch from remote(Product) repo.

For Project delivery, we added upstream for this repo.

<https://stefanbauer.me/articles/how-to-keep-your-git-fork-up-to-date>

#### How to add an upstream

1. Add remote upstream for this repo,

```bash

$ git remote add upstream git@bitbucket.org:example/example.git

```

2. Check your upstream is available for the repo

```bash

$ git remote -v

origin git@bitbucket.org:example/spark-example.git (fetch)

origin	git@bitbucket.org:example/spark-example.git (push)

upstream	git@bitbucket.org:example/example.git (fetch)

upstream	git@bitbucket.org:example/example.git (push)

```

3. Get upto date with upstream

```bash

$ git fetch upstream

```

4. Merge Local master branch with Upstream master branch,

```bash

$ git merge upstream/master master

```

#### Note: You will be always rebase the upstream release/master branch with local release branch

Always rebase with the master branch,

```bash

$ git checkout  release/2020_Q4_3.1

$ git rebase upstream/release/2020/Q4

```
#### Create Release/Tag

Add this alias in `~/.gitconfig`

```bash

addtag = "!f() { git tag -l | xargs git tag -d && git fetch -p && git tag $1; } ; f"

ptags = push origin --tags

dtags = "!f() { git tag -d $1 && git push origin :refs/tags/$1; } ; f"

```

`git addtag v3.1.67`

`git ptags v3.1.67`

### Note: Keep the commit message starting with keyword to move across release

Example:

```bash

git commit -m "Customer: commit the changes"

```
