from flask import Flask, request,render_template, jsonify, Response
import random
import glob
import os
import uuid
import urllib2
import time
 
app = Flask(__name__)
# print "hello"
@app.route("/upload_phrase", methods=['POST'])
def classify_url():
    if request.method == 'POST':

        phrase = request.form['phrase']
        phrase = parse_and_dump.get_language_text(phrase)

        counter = len(phrase)

        test_splice_text.wrapper_main(phrase)
        # time.sleep(counter / 2)

        print phrase
        return jsonify({'phrase':phrase})

    else:
        #get 10 most similar and return
        return 

@app.route("/")
def hello():
    print "hi"
    return render_template('index.html')

def gen(camera):
    while True:
        frame = camera.get_frame()
        yield (b'--frame\r\n'
               b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n\r\n')

@app.route('/video_feed')
def video_feed():
    return Response(gen(VideoCamera()),
                    mimetype='multipart/x-mixed-replace; boundary=frame')

if __name__ == "__main__":
    app.run(host = '0.0.0.0', debug = True)
