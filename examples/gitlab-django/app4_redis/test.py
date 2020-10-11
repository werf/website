import requests

def batch_mailing():

    return requests.post(

        "https://api.mailgun.net/v3/sandboxdad71610c342445aa4ab9bd5fe448eaf.mailgun.org/messages",
        auth=("api", "b7298d7dd921c92c1ec29cf0e499513c-0d2e38f7-00870ddd"),
        data = {
            "from": "noreply@sandboxdad71610c342445aa4ab9bd5fe448eaf.mailgun.org",
            "to":["rodionovs45qwe@gmail.com"],
            "subject":"New message",
            "text": "Hi"
        })

# Run the function
print(batch_mailing())