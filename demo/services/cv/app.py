from flask import Flask, request, jsonify
app = Flask(__name__)

@app.route("/health")
def health():
    return jsonify({"status":"ok"})

@app.route("/detect", methods=["POST"])
def detect():
    # pretend we processed image
    return jsonify({
        "detections": [
            {"class":"box","confidence":0.95,"bbox":[10,20,200,150],"ocr": "SKU-1234"}
        ],
        "width":1280,"height":720
    })

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8001)
