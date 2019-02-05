workflow "Splash PRs" {
  on = "pull_request"
  resolves = ["calavera/splashed-prs@master"]
}

action "calavera/splashed-prs@master" {
  uses = "docker://calavera/splashed-prs:latest"
  secrets = ["GITHUB_TOKEN", "UNSPLASH_CLIENT_ID"]
  env = {
    UNSPLASH_QUERY = "cute animal"
    UNSPLASH_ORIENTATION = "portrait"
    DEBUG = "true"
  }
}
