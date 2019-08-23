// @flow weak

import React, { Component } from 'react'


export class AddIcon extends Component {
    render() {
        return (
            <svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                <path d="M0 0h24v24H0z" fill="none" />
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm5 11h-4v4h-2v-4H7v-2h4V7h2v4h4v2z" />
            </svg>
        )
    }
}

export class EditIcon extends Component {
    render() {
        return (
            <svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z" />
                <path d="M0 0h24v24H0z" fill="none" />
            </svg>
        )
    }
}

export class RemoveIcon extends Component {
    render() {
        return (
            <svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z" />
                <path d="M0 0h24v24H0z" fill="none" />
            </svg>
        )
    }
}

export class OpenInNewIcon extends Component {
    render() {
        return (
            <svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                <path d="M0 0h24v24H0z" fill="none" />
                <path d="M19 19H5V5h7V3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.1 0 2-.9 2-2v-7h-2v7zM14 3v2h3.59l-9.83 9.83 1.41 1.41L19 6.41V10h2V3h-7z" />
            </svg>
        )
    }
}

export class ArrowDownIcon extends Component {
    render() {
        return (
            <svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                <path d="M7.41 7.84L12 12.42l4.59-4.58L18 9.25l-6 6-6-6z" />
                <path d="M0-.75h24v24H0z" fill="none" />
            </svg>
        )
    }
}

export class ArrowRightIcon extends Component {
    render() {
        return (
            <svg fill="#000000" height="24" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg">
                <path d="M8.59 16.34l4.58-4.59-4.58-4.59L10 5.75l6 6-6 6z" />
                <path d="M0-.25h24v24H0z" fill="none" />
            </svg>
        )
    }
}

export class HamburgerIcon extends Component {
    render() {
        return (
            <svg fill="#000000" viewBox="0 0 24 24" height="24" width="24" xmlns="http://www.w3.org/2000/svg">
                <path d="M3 13h18v-2H3v2zm0 6h18v-2H3v2zM3 5v2h18V5H3z" />
                <path d="M0 0h24v24H0V0z" fill="none" />
            </svg>
        )
    }
}