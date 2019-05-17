// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.


import Client from '../client';

export function startUp() {
    return async (dispatch, getState) => {
        try {
            await Client.startUp();
        } catch (error) {
            return {error};
        }
        return {data: true};
    };
}
