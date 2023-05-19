pwd
echo $WERF_DIR
werf converge --repo k3d-registry.sample-app.test:5000/sample-app
curl sample-app.test/ping
