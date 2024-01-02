import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    stages: [
        {duration: '10s' , target: 100}, // below normal load
        {duration: '30s' , target: 100},
        {duration: '10s' , target: 1400}, // spike to 1400 users
        {duration: '1m' , target: 1400}, // stay at 1400 for 3 minutes
        {duration: '10s' , target: 100}, // scale down. Recovery stage.
        {duration: '1m' , target: 100},
        {duration: '10s' , target: 0},
    ]
};

// Function to generate a random string
function generateRandomString(length) {
    const characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    let result = '';
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    return result;
}

// Function to generate a random timestamp
function generateRandomTimestamp() {
    const date = new Date(Date.now() - Math.floor(Math.random() * 365 * 24 * 60 * 60 * 1000));
    return date.toISOString();
}

// Function to generate a random payload
function generateRandomPayload() {
    return {
        level: ['error', 'info', 'debug', 'fatal'][Math.floor(Math.random() * 4)],
        message: `Log message: ${generateRandomString(10)}`,
        resourceId: `server-${Math.floor(Math.random() * 9000) + 1000}`,
        timestamp: generateRandomTimestamp(),
        traceId: `${generateRandomString(3)}-${generateRandomString(3)}-${Math.floor(Math.random() * 900) + 100}`,
        spanId: `span-${Math.floor(Math.random() * 900) + 100}`,
        commit: generateRandomString(7),
        metadata: {
            parentResourceId: `server-${Math.floor(Math.random() * 9000) + 1000}`,
        },
    };
}


export default () => {
    let url = 'http://localhost:1323/public/ingest-kafka';

    let payload = JSON.stringify(generateRandomPayload());

    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.post(url, payload, params);

    // Check if the request was successful (status code 2xx)
    check(res, {
        'is status 2xx': (r) => r.status >= 200 && r.status < 300,
    });

    // Simulate a delay between requests
    sleep(0.1);  // 100 milliseconds
}
