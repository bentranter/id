Package psa should satisfy a few requirements:

  1. Be passwordless
  2. Be sessionless
  3. Be incredibly easy to create new providers/stratgies
     for.
  4. Work for any middleware setup/config.
  5. Work with either cookies localstorage.
  6. Be easy to set up continuous integration tests for.
  7. Have good support for every Oauth feature.
  8. Which means actual support for refresh tokens.

Let's figure out our dreamcode.

```go
http.HandleFunc("/auth/github/authorize", GitHub.Authorize)
http.HandleFunc("/auth/github/callback", GitHub.Callback)
```

Probably gonna need context package to pass around a few important things.

IDEA: For the three-legged Oauth stuff, pass a context over a channel that contains the config. The config can be updated with each value as it becomes satisfied?

You can take the same approach for CSRF tokens as well? Idk if this works or can really be called sessionless... well the idea I supposed is that if you shut your server down, you can restart without having to worry about users' sessions being lost, or you can add another process??? and the sessions will still work. I think my idea is okay.
