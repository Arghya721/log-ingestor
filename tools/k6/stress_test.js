/* 
Stress Testing is a type of performance testing that validates the highest limit of your application’s performance in terms of scalability, stability, and reliability. It is performed to determine how the system under test will behave when the load on the system is raised above the expected maximum.

Run a stress test to:
    - Determine how the system behaves under extreme load
    - Determine the maximum capacity of the system
    - Determine the breaking point of the system and its failure mode
    - Determine if your system recovers without manual intervention after the stress test
*/

import http from 'k6/http';
import { check, sleep } from 'k6';

// Export options object
export let options = {
    // Skip TLS verification
    insecureSkipTLSVerify: true,
  
    // Do not reuse connections
    noConnectionReuse: false,
  
    // Stages of the test
    stages: [
      // Below normal load
      { duration: '2m', target: 100 },
      { duration: '5m', target: 100 },
  
      // Normal load
      { duration: '2m', target: 200 },
      { duration: '5m', target: 200 },
  
      // Around the breaking point
      { duration: '2m', target: 300 },
      { duration: '5m', target: 300 },
  
      // Beyond the breaking point
      { duration: '2m', target: 400 },
      { duration: '5m', target: 400 },
  
      // Recovery stage
      { duration: '10m', target: 0 },
    ],
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
    let url = 'http://localhost:1323/public/ingest';

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
