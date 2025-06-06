![wera-logo](wera-logo-white.png)

# Wera

[Project Info](https://eliasalerno.ch/work/wera)

Cities are heating up during the summer, and deaths related to this are on the rise. To mitigate this ever-increasing risk, we have prototyped, built, and user-tested an entirely open-source, personalized city climate warning system called Wera. Wera can access the high-density temperature sensor network in Zurich, along with meteorological data and the user's location, to formulate a weather and risk potential report with Llama3, which is then synthesized into either German or English for the end user. Wera primarily targets seniors, as they are the most affected by the city's hotspots. It features a very intuitive and accessible interface, making it a reliable weather risk informer and a helpful weather assistant. Wera comes in two form factors: the Kub, its strengths being its constant visibility, personal feeling, and accessibility, and the Fong, which is widely deployable and offers the option to ask back questions directly on your phone.

[Youtube Video](https://www.youtube.com/watch?v=xea25izHtzc)

## Table of Contents

- [Wera](#wera)
  - [Table of Contents](#table-of-contents)
  - [Hardware](#hardware)
    - [Kub](#kub)
    - [API's](#apis)
  - [Server Quickstart:](#server-quickstart)
    - [**Build the server (with docker):**](#build-the-server-with-docker)
    - [**Build and run the the server (without Docker):**](#build-and-run-the-the-server-without-docker)
  - [Radio (Raspberry PI 4) Quickstart](#radio-raspberry-pi-4-quickstart)
    - [Fetch dependencies and build:](#fetch-dependencies-and-build)
- [Docs](#docs)
  - [Endpoints](#endpoints)
    - [**GET /ok**](#get-ok)
    - [**GET /cityclimate**](#get-cityclimate)
    - [**POST /cityclimategps**](#post-cityclimategps)
    - [**GET /meteoblue**](#get-meteoblue)
    - [**POST /weathergps**](#post-weathergps)
    - [**POST /weather**](#post-weather)
    - [**GET /hotareas**](#get-hotareas)
    - [**POST /hotareasgps**](#post-hotareasgps)
    - [**POST /speak**](#post-speak)
  - [Dependencies](#dependencies)
    - [Server](#server)
    - [Kub (Raspberry Pi 4)](#kub-raspberry-pi-4)
    - [License](#license)

---

## Hardware

### Kub

![kub-schematics](kub/kub_schematics/kub_schematics.png)

- Raspberry Pi (we used the Pi 4 model b)
- 1x GPS Module
- 2x Status LED's (1x Green 1x Red)
- 1x Rotary encoder
- 1x Speaker, easiest if it has it's own amplifier built in

---

### API's

- [MeteoBlue](https://www.meteoblue.com/de/weather-api/index/overview)
- [VoiceRSS API (backup TTS)](https://www.voicerss.org/personel/)
- [Groq](https://console.groq.com/docs/quickstart)
- [Unreal TTS (backup TTS)](https://unrealspeech.com/onboard)
- [Google TTS](https://cloud.google.com/text-to-speech?hl=de)

## Server Quickstart:

Be sure to get all API Keys, .env is formatted like this in the `/server` directory

```
METEO_API_KEY=key
TTS_API_KEY=key // -> VoiceRSS API (backup TTS), package can be removed also
GROQ_API_KEY=key
UNREAL_SPEECH_API_KEY=key // -> Unreal TTS (backup TTS), package can be removed also
GOOGLE_APPLICATION_CREDENTIALS=key.json-file
```

### **Build the server (with docker):**

```bash
cd server
docker build -t wera-server .
docker run -d -p 8080:8080 wera-server
```

### **Build and run the the server (without Docker):**

```bash
cd server
go mod tidy
go build cmd/server/main.go
./main
```

## Radio (Raspberry PI 4) Quickstart

### Fetch dependencies and build:

1. **Navigate to the kub directory:**

```bash
cd kub
```

2. **Initialize the Go module** (if not already done):

```bash
go mod tidy
```

3. **Fetch dependencies:**

```bash
go get github.com/stianeikeland/go-rpio/v4 go.bug.st/serial.v1
```

4. **Build the kub application:**

```bash
go build kub.go
```

5. **Run the built application:**

```bash
./kub
```

6. **Create a Service for the Kub**

```bash
sudo nano /etc/systemd/system/kub.service
```

7. **Add the following content** to the service file. Adjust the `ExecStart` path to point to your executable and modify other settings as necessary:

```ini
[Unit]
Description=Kub Service
After=network.target

[Service]
ExecStart=/home/user/wera/kub/kub
WorkingDirectory=/home/user/wera/kub
StandardOutput=inherit
StandardError=inherit
Restart=always
User=root
RestartSec=10
StartLimitIntervalSec=500
StartLimitBurst=5


[Install]
WantedBy=multi-user.target
```

8. **Reload the systemd manager configuration** to read the newly created service file:

```bash
sudo systemctl daemon-reload
```

9. **Enable the service** to start on boot:

```bash
sudo systemctl enable goradio.service
```

10. **Start the service** immediately to test it:

```bash
sudo systemctl start goradio.service
```

11. **Check the status** of the service to ensure it is active and running:

```bash
sudo systemctl status goradio.service
```

---

# Docs

Render Hosted Service (2min startup time): https://spatial-interaction.onrender.com

## Endpoints

### **GET /ok**

Responds with OK status 200.

### **GET /cityclimate**

Responds with the sensor dataset of the ZHAW Grid, currently only has access to about 50 sensors and the temperature data only.

### **POST /cityclimategps**

Responds with the sensor dataset of the ZHAW Grid, sorted by distance (closest to furthest).

Sample Request Body:

```JSON
{
"device_id": "username",
"Latitude": 47.3653466,
"Longitude": 8.5282651,
"Language": "english"
}
```

Sample Response Body:

```JSON
[
  {
    "id": "03400120",
    "name": "Sihlhölzlistrasse",
    "timestamp": "1716555600",
    "values": 22.8375,
    "colors": "#9baf33",
    "active": 1,
    "Distance": 0.3995830946586305,
    "geometry": {
      "type": "point",
      "coordinates": {
        "lon": 8.528826,
        "lat": 47.36892
      }
    }
  }
]
```

### **GET /meteoblue**

Responds with a 3-hour forecast from Meteoblue data, also provides a 24-hour overview. Data contains temperature, wind, rain, and some more. Each request takes 8000 tokens and our free API is limited to 10M, so please make only as many requests as needed.

### **POST /weathergps**

Fetches and combines data from both the MeteoBlue and CityClimate APIs, processes it to compute different data points, to create a current weather analysis and forecast. It uses the GPS data of the kub to deliver data for the current location of the radio, constructs a nice response, and generates a speech file (.wav) which is returned as an audio stream.

Sample Request Body:

```JSON
{
  "device_id": "Device_1",
  "Latitude": 47.3653466,
  "Longitude": 8.5282651,
  "Language": "english"
}
```

Sample text it generates and speaks:

```plaintext
"Good morning! It's a lovely day outside! The temperature has been quite pleasant, a gentle breeze is blowing and the temperature is just right - not too hot, not too cold. Just perfect. Make sure to stay hydrated and take a break if you're spending time outdoors. And remember, on especially warm days, please be extra careful to avoid heat exhaustion. Stay cool and comfortable!"
```

### **POST /weather**

Fetches and combines data from both the MeteoBlue and CityClimate APIs, processes it to compute different data points, to create a current weather analysis and forecast, constructs a descriptive sentence, and generates a speech file (.wav) which is returned as an audio stream. This endpoint is used as a relay in case weathergps fails.

Sample Request Body:

```JSON
{
"device_id": "username",
"Language": "english"
}
```

The response contains:

- The current average temperature of the sensor grid.
- The temperature and wind speed according to MeteoBlue.

Sample text it generates and speaks:

```plaintext
The current average temperature of the Sensor Grid is 22.50 degrees Celsius. According to MeteoBlue, the temperature is 20.10 degrees Celsius with a windspeed of 3.5 meters per second.
```

### **GET /hotareas**

Resonds with an array of sensors that are over a certain set threshold, on default it is set to 28 degrees celsius.

Sample Response Body:

```json
[
  {
    "id": "0340011C",
    "name": "Zollikerstrasse",
    "timestamp": "1716555600",
    "values": 28.8025,
    "colors": "#743933",
    "active": 1,
    "geometry": {
      "type": "point",
      "coordinates": {
        "lon": 8.556188,
        "lat": 47.3617
      }
    }
  }
]
```

### **POST /hotareasgps**

Resonds with an array of sensors that are over a certain set threshold, on default it is set to 28 degrees celsius, sorted by distance to the post location.

Sample Request Body:

```JSON
{
  "device_id": "Device_1",
  "Latitude": 47.3653466,
  "Longitude": 8.5282651
}
```

Sample Response Body:

```json
[
  {
    "id": "0340011C",
    "name": "Zollikerstrasse",
    "timestamp": "1716555600",
    "values": 28.8025,
    "colors": "#743933",
    "active": 1,
    "geometry": {
      "type": "point",
      "coordinates": {
        "lon": 8.556188,
        "lat": 47.3617
      }
    }
  }
]
```

### **POST /speak**

Text to speech endpoint. (Unreal speech)

Sample Request Body:

```JSON
{
  "text": "Note: As soon as the green light comes on, I can provide more detailed information, this normally takes 15 minutes. The kub can be used nevertheless ",
  "language": "english"
}

```

Sample Response Body:

```JSON
{
  "file": "output.wav",
  "message": "Speech generated successfully"
}
```

---

## Dependencies

### Server

- `github.com/jackc/pgx/v4` for PostgreSQL database interaction.
- `github.com/gorilla/mux` for routing.
- `google.golang.org/genproto/googleapis/cloud/texttospeech/v1` & `cloud.google.com/go/texttospeech/apiv1` for Google Text to Speech integration

### Kub (Raspberry Pi 4)

- `github.com/stianeikeland/go-rpio/v4`for GPIO pin support.
- `go.bug.st/serial.v1`for Serial Support.

### License

This Project lies under the MIT License.
