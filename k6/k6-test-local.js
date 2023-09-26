import http from 'k6/http';
import {sleep} from 'k6';

export const options = {
    vus: 80,
    duration: '600s',
    noConnectionReuse: true,
};

export default function () {
    http.get('http://localhost:8080/ping');
    sleep(0.01);
}


