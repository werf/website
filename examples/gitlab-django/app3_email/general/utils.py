import requests

def batch_mailing(email):

    return requests.post(

        "https://api.mailgun.net/v3/sandboxdad71610c342445aa4ab9bd5fe448eaf.mailgun.org/messages",
        auth=("api", ""),
        data = {
            "from": "noreply@sandboxdad71610c342445aa4ab9bd5fe448eaf.mailgun.org",
            "to":[email],
            "subject":"New message",
            "text": "Hi"
        })

