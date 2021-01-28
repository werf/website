Run the following port forwarding command in the terminal:

{% raw %}

sudo apt-get install -y socat
socat -d -d TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000
{% endraw %}
