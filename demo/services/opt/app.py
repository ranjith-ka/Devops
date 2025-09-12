from flask import Flask, request, jsonify
app = Flask(__name__)

@app.route("/health")
def health():
    return jsonify({"status":"ok"})

@app.route("/optimize", methods=["POST"])
def optimize():
    return jsonify({"route":"demo-route","distance_km":12.5})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8003)
