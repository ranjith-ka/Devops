from flask import Flask, request, jsonify
app = Flask(__name__)

@app.route("/health")
def health():
    return jsonify({"status":"ok"})

@app.route("/classify", methods=["POST"])
def classify():
    text = request.json.get("text","")
    if "customs" in text.lower():
        label = "customs_issue"
    elif "delay" in text.lower():
        label = "delay"
    else:
        label = "general"
    return jsonify({"label": label, "confidence": 0.9})

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000)
