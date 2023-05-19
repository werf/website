helm upgrade --wait --install --create-namespace --namespace sample-app sample-app .helm/ --set werf.image.app=${APP_IMAGE} --set host=sample-app.test --timeout 60s
# timeout error occurred, let's dig into the cluster to get an actual error
kubectl -n sample-app get pod
# there is crashed pod
kubectl -n sample-app describe pod $(kubectl -n sample-app get pod | grep Crash | cut -d' ' -f1)
# check out pod logs
kubectl -n sample-app logs $(kubectl -n sample-app get pod | grep Crash | cut -d' ' -f1)
# let's check out deployment manifest
cat .helm/templates/deployment.yaml
# typo in the command field, let's fix it
sed -e 's|"star"|"start"|' -i .helm/templates/deployment.yaml
helm upgrade --wait --install --create-namespace --namespace sample-app sample-app .helm/ --set werf.image.app=${APP_IMAGE} --set host=sample-app.test --timeout 60s
curl sample-app.test/ping
