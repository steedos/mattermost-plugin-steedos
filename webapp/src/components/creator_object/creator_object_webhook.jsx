// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';
import PropTypes from 'prop-types';

import { makeStyleFromTheme } from 'mattermost-redux/utils/theme_utils';


export default class PostTypeWorkflowWebhook extends React.PureComponent {
    static propTypes = {

        /*
         * The post to render the message for.
         */
        post: PropTypes.object.isRequired,

        /**
         * Set to render post body compactly.
         */
        compactDisplay: PropTypes.bool,

        /**
         * Flags if the post_message_view is for the RHS (Reply).
         */
        isRHS: PropTypes.bool,

        /*
         * Logged in user's theme.
         */
        theme: PropTypes.object.isRequired,
    };

    static defaultProps = {
        mentionKeys: [],
        compactDisplay: false,
        isRHS: false,
    };

    constructor(props) {
        super(props);

        this.state = {
        };
    }

    render() {
        const style = getStyle(this.props.theme);
        const post = this.props.post;
        const props = post.props || {};

        let preText = '';
        if ('create' == props.action) {
            preText = `${props.actionUserInfo.name}新建了${props.objectDisplayName}: `;
        } else if ('update' == props.action) {
            preText = `${props.actionUserInfo.name}更新了${props.objectDisplayName}: `;
        } else if ('delete' == props.action) {
            preText = `${props.actionUserInfo.name}删除了${props.objectDisplayName}: `;
        }

        let title = props.data[props.nameFieldKey];

        return (
            <div>
                {preText}
                <div style={style.attachment}>
                    <div style={style.content}>
                        <div style={style.container}>
                            <b>{'标题'}</b>
                            <br />
                            <a
                                target='_blank'
                                href={props.redirectUrl}
                            >
                                {title}
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        attachment: {
            marginLeft: '-20px',
            position: 'relative',
        },
        content: {
            borderRadius: '4px',
            borderStyle: 'solid',
            borderWidth: '0px',
            borderColor: '#BDBDBF',
            margin: '5px 0 5px 20px',
            padding: '2px 5px',
        },
        container: {
            borderLeftStyle: 'solid',
            borderLeftWidth: '4px',
            paddingTop: '5px',
            paddingTight: '10px',
            paddingBottom: '5px',
            paddingLeft: '10px',
            borderLeftColor: '#89AECB',
        },
        body: {
            overflowX: 'auto',
            overflowY: 'hidden',
            paddingRight: '5px',
            width: '100%',
        },
    };
});
