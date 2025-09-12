from flask import Flask, jsonify
app = Flask(__name__)

@app.route("/health")
def health():
    return jsonify({"status":"ok"})

@app.route("/forecast", methods=["GET","POST"])
def forecast():
    return jsonify({"forecast":[{"ts":"2025-09-13","pred":0.2}]})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8002)
