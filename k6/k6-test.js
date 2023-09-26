import http from 'k6/http';
import {sleep} from 'k6';

export const options = {
    vus: 300,
    duration: '120s',
};

export default function () {
    http.get('http://pingserver:8080/ping');
    sleep(1);
}