## Development

Build and run locally:
```shell
echo "ru.localhost" | sudo tee -a /etc/hosts
werf compose up --env development --dev --docker-compose-command-options='-d'
```

Now access the site:
```shell
curl ru.localhost
```
