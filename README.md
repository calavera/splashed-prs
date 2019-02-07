# Splashed Pull Requests

GitHub action that adds pictures to the description of every Pull Requests.

[Unsplash](https://unsplash.com/) is a community to discover pictures gifted by amazing photographers.

## Usage

This action requires you to register a new Unsplash application 
https://unsplash.com/developers before using it.

### Secrets

- UNSPLASH_CLIENT_ID: Your API Client ID from Unsplash. You don't need the 
  client secret because the endpoint that this action uses is public.

### Environment variables

- UNSPLASH_QUERY: Filter random photo with a search term.
- UNSPLASH_ORIENTATION: Select photo orientation.
- DEBUG: Print more information during the action execution.

## Disclaimer

This kind of usage is against Unsplash guidelines, but it was fun to build, and 
you can use it in developer mode if the volume of requests is very low:

https://help.unsplash.com/api-guidelines/unsplash-api-guidelines

> The API is to be used for non-automated, high-quality, and authentic experiences. 
