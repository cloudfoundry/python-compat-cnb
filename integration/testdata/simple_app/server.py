from flask import Flask
import platform
app = Flask(__name__)

@app.route('/')
def hello_world():
    return 'Hello, World!'

@app.route('/version')
def version():
    version = platform.python_version()
    return "Python version: " + version

if __name__ == "__main__":
    app.run()
