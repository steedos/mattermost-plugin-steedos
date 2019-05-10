import {id as pluginId} from './manifest';
import React from 'react';

// export default class Plugin {
//     // eslint-disable-next-line no-unused-vars
//     initialize(registry, store) {
//         // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
//     }
// }

// window.registerPlugin(pluginId, new Plugin());


// Courtesy of https://feathericons.com/
const Icon = () => <i className='icon fa fa-plug'/>;

export default class Plugin {
    initialize(registry, store) {
        registry.registerChannelHeaderButtonAction(
            // icon - JSX element to use as the button's icon
            <Icon />,
            // action - a function called when the button is clicked, passed the channel and channel member as arguments
            // null,
            () => {
                alert("Hello World!");
            },
            // dropdown_text - string or JSX element shown for the dropdown button description
            "Hello World",
        );
    }
}

window.registerPlugin(pluginId, new Plugin());
