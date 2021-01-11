from flask import Flask
app = Flask(__name__)

@app.route('/pyapp')
def hello_world():
    return 'Hello Python World!'

if __name__ == '__main__':
    app.run()