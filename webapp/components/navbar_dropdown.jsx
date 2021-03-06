// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import $ from 'jquery';
import ReactDOM from 'react-dom';
import * as Utils from 'utils/utils.jsx';
import * as GlobalActions from 'action_creators/global_actions.jsx';

import AboutBuildModal from './about_build_modal.jsx';
import TeamMembersModal from './team_members_modal.jsx';
import ToggleModalButton from './toggle_modal_button.jsx';
import UserSettingsModal from './user_settings/user_settings_modal.jsx';

import Constants from 'utils/constants.jsx';

import {FormattedMessage} from 'react-intl';
import {Link} from 'react-router';

import React from 'react';

export default class NavbarDropdown extends React.Component {
    constructor(props) {
        super(props);
        this.blockToggle = false;

        this.handleAboutModal = this.handleAboutModal.bind(this);
        this.aboutModalDismissed = this.aboutModalDismissed.bind(this);

        this.state = {
            showUserSettingsModal: false,
            showAboutModal: false
        };
    }
    handleAboutModal() {
        this.setState({showAboutModal: true});
    }
    aboutModalDismissed() {
        this.setState({showAboutModal: false});
    }
    componentDidMount() {
        $(ReactDOM.findDOMNode(this.refs.dropdown)).on('hide.bs.dropdown', () => {
            $('.sidebar--left .dropdown-menu').scrollTop(0);
            this.blockToggle = true;
            setTimeout(() => {
                this.blockToggle = false;
            }, 100);
        });
    }
    componentWillUnmount() {
        $(ReactDOM.findDOMNode(this.refs.dropdown)).off('hide.bs.dropdown');
    }
    render() {
        var teamLink = '';
        var inviteLink = '';
        var manageLink = '';
        var sysAdminLink = '';
        var adminDivider = '';
        var currentUser = this.props.currentUser;
        var isAdmin = false;
        var isSystemAdmin = false;
        var teamSettings = null;
        let integrationsLink = null;

        if (currentUser != null) {
            isAdmin = Utils.isAdmin(currentUser.roles);
            isSystemAdmin = Utils.isSystemAdmin(currentUser.roles);

            inviteLink = (
                <li>
                    <a
                        href='#'
                        onClick={GlobalActions.showInviteMemberModal}
                    >
                        <FormattedMessage
                            id='navbar_dropdown.inviteMember'
                            defaultMessage='Invite New Member'
                        />
                    </a>
                </li>
            );

            if (this.props.teamType === Constants.OPEN_TEAM && global.window.mm_config.EnableUserCreation === 'true') {
                teamLink = (
                    <li>
                        <a
                            href='#'
                            onClick={GlobalActions.showGetTeamInviteLinkModal}
                        >
                            <FormattedMessage
                                id='navbar_dropdown.teamLink'
                                defaultMessage='Get Team Invite Link'
                            />
                        </a>
                    </li>
                );
            }
        }

        if (isAdmin) {
            manageLink = (
                <li>
                    <ToggleModalButton dialogType={TeamMembersModal}>
                        <FormattedMessage
                            id='navbar_dropdown.manageMembers'
                            defaultMessage='Manage Members'
                        />
                    </ToggleModalButton>
                </li>
            );

            adminDivider = (<li className='divider'></li>);

            teamSettings = (
                <li>
                    <a
                        href='#'
                        data-toggle='modal'
                        data-target='#team_settings'
                    >
                        <FormattedMessage
                            id='navbar_dropdown.teamSettings'
                            defaultMessage='Team Settings'
                        />
                    </a>
                </li>
            );
        }

        if (window.mm_config.EnableIncomingWebhooks === 'true' || window.mm_config.EnableOutgoingWebhooks === 'true') {
            if (isAdmin || window.EnableAdminOnlyIntegrations !== 'true') {
                integrationsLink = (
                    <li>
                        <Link to={'/settings/integrations'}>
                            <FormattedMessage
                                id='navbar_dropdown.integrations'
                                defaultMessage='Integrations'
                            />
                        </Link>
                    </li>
                );
            }
        }

        if (isSystemAdmin) {
            sysAdminLink = (
                <li>
                    <Link
                        to={'/admin_console'}
                    >
                        <FormattedMessage
                            id='navbar_dropdown.console'
                            defaultMessage='System Console'
                        />
                    </Link>
                </li>
            );
        }

        var teams = [];

        if (global.window.mm_config.EnableTeamCreation === 'true') {
            teams.push(
                <li key='newTeam_li'>
                    <Link
                        key='newTeam_a'
                        target='_blank'
                        to={Utils.getWindowLocationOrigin() + '/signup_team'}
                    >
                        <FormattedMessage
                            id='navbar_dropdown.create'
                            defaultMessage='Create a New Team'
                        />
                    </Link>
                </li>
            );
        }

        let helpLink = null;
        if (global.window.mm_config.HelpLink) {
            helpLink = (
                <li>
                    <Link
                        target='_blank'
                        to={global.window.mm_config.HelpLink}
                    >
                        <FormattedMessage
                            id='navbar_dropdown.help'
                            defaultMessage='Help'
                        />
                    </Link>
                </li>
            );
        }

        let reportLink = null;
        if (global.window.mm_config.ReportAProblemLink) {
            reportLink = (
                <li>
                    <Link
                        target='_blank'
                        to={global.window.mm_config.ReportAProblemLink}
                    >
                        <FormattedMessage
                            id='navbar_dropdown.report'
                            defaultMessage='Report a Problem'
                        />
                    </Link>
                </li>
            );
        }

        return (
            <ul className='nav navbar-nav navbar-right'>
                <li
                    ref='dropdown'
                    className='dropdown'
                >
                    <a
                        href='#'
                        className='dropdown-toggle'
                        data-toggle='dropdown'
                        role='button'
                        aria-expanded='false'
                    >
                        <span
                            className='dropdown__icon'
                            dangerouslySetInnerHTML={{__html: Constants.MENU_ICON}}
                        />
                    </a>
                    <ul
                        className='dropdown-menu'
                        role='menu'
                    >
                        <li>
                            <a
                                href='#'
                                onClick={() => this.setState({showUserSettingsModal: true})}
                            >
                                <FormattedMessage
                                    id='navbar_dropdown.accountSettings'
                                    defaultMessage='Account Settings'
                                />
                            </a>
                        </li>
                        {inviteLink}
                        {teamLink}
                        <li>
                            <Link to={'/' + this.props.teamName + '/logout'}>
                                <FormattedMessage
                                    id='navbar_dropdown.logout'
                                    defaultMessage='Logout'
                                />
                            </Link>
                        </li>
                        {adminDivider}
                        {teamSettings}
                        {integrationsLink}
                        {manageLink}
                        {sysAdminLink}
                        {teams}
                        <li className='divider'></li>
                        {helpLink}
                        {reportLink}
                        <li>
                            <a
                                href='#'
                                onClick={this.handleAboutModal}
                            >
                                <FormattedMessage
                                    id='navbar_dropdown.about'
                                    defaultMessage='About Mattermost'
                                />
                            </a>
                        </li>
                        <UserSettingsModal
                            show={this.state.showUserSettingsModal}
                            onModalDismissed={() => this.setState({showUserSettingsModal: false})}
                        />
                        <AboutBuildModal
                            show={this.state.showAboutModal}
                            onModalDismissed={this.aboutModalDismissed}
                        />
                    </ul>
                </li>
            </ul>
        );
    }
}

NavbarDropdown.defaultProps = {
    teamType: ''
};
NavbarDropdown.propTypes = {
    teamType: React.PropTypes.string,
    teamDisplayName: React.PropTypes.string,
    teamName: React.PropTypes.string,
    currentUser: React.PropTypes.object
};
