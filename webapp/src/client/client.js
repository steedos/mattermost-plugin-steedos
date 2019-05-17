// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import request from 'superagent';

import {id} from '../manifest';

const {Client4} = require('mattermost-redux/client');

export default class Client {
    constructor() {
        this.url = '/plugins/' + id;
    }

    startUp = async () => {
        return this.doGet(`${this.url}/startup`);
    }

    doGet = async (url, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';

        console.log('Client4.getToken(): ', Client4.getToken())

        try {
            const response = await request.
                get(url).
                set(headers).
                type('application/json').
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }
}
