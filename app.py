from flask import Flask, request,render_template, jsonify, Response
import random
import glob
import os
import uuid
import urllib2
import time
import httplib, urllib
import requests
import json
 
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
def index_main():
    print "rendering website"
    return render_template('index.html', name = "hahahahahahah")

@app.route("/add_activity",methods=["POST"])
def add_activity(req = None):
    print "rendering post activity"
    # print req
    # print request.form["description"]
    # print request

    name = request.form['activity_name']
    print name, "  activity_name"
    description = request.form['description']
    print description , " desc"
    try:
        capacity = int(request.form['capacity'])
    except:
        capacity = 12
    print capacity, "capacity"
    location = request.form['location']
    print location , "location"
    x = request.form['loclat']
    print x, "loclat"
    y = request.form['loclong']
    print y, "locLong"
    point = str(x) + "," + str(y)
    print point, "point"
    start_time = request.form['start_time']
    end_time = request.form['end_time']
    owner = 555

    
    data_r = {
        "name" : name,
        "desc" : description,
        "cap" : capacity,
        "loc" : location,
        "point" : point,
        "start" : start_time,
        "end" : end_time,
        "owner" : owner


    }

    data_r_json = json.dumps(data_r)

    r = requests.post("http://10.10.200.66:8080/activity", data= data_r_json)
    print(r.status_code, r.reason)

    return render_template('submit_form.html', name = "add_activity_meh")

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
