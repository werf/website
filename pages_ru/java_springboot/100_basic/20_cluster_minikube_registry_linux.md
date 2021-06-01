Запустите следующий проброс портов в отдельном терминале:

{% raw %}
```shell
sudo apt-get install -y socat
socat -d -d TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000
```
{% endraw %}