import {id as pluginId} from './manifest';

import PostTypeWorkflowWebhook from './components/workflow';
import CreatorObjectWebhook from './components/creator_object';

// export default class Plugin {
//     // eslint-disable-next-line no-unused-vars
//     initialize(registry, store) {
//         // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
//     }
// }

// window.registerPlugin(pluginId, new Plugin());


// Courtesy of https://feathericons.com/
// import {startUp} from './actions';

// const Icon = () => <i className='icon fa fa-plug'/>;

export default class Plugin {
    initialize(registry, store) {
        // registry.registerChannelHeaderButtonAction(
        //     // icon - JSX element to use as the button's icon
        //     <Icon />,
        //     // action - a function called when the button is clicked, passed the channel and channel member as arguments
        //     // null,
        //     (channel) => {
        //         startUp()();
        //     },
        //     // dropdown_text - string or JSX element shown for the dropdown button description
        //     "start up!",
        // );
        registry.registerPostTypeComponent('custom_workflow_webhook', PostTypeWorkflowWebhook);
        registry.registerPostTypeComponent('custom_object_webhook', CreatorObjectWebhook);
    }
}

window.registerPlugin(pluginId, new Plugin());
