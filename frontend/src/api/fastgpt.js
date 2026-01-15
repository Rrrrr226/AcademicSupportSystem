import axios from 'axios';
import { BASE_URL } from './config';

const getAuthHeader = (token) => ({ Authorization: `Bearer ${token}` });

export const getAppList = (token, page = 1, pageSize = 10) => {
    return axios.post(`${BASE_URL}/fastgpt/apps/list`, {
        offset: (page - 1) * pageSize,
        limit: pageSize
    }, {
        headers: getAuthHeader(token)
    });
};

export const createApp = (token, data) => {
    return axios.post(`${BASE_URL}/fastgpt/apps/create`, data, {
        headers: getAuthHeader(token)
    });
};

export const updateApp = (token, data) => {
    return axios.post(`${BASE_URL}/fastgpt/apps/update`, data, {
        headers: getAuthHeader(token)
    });
};

export const deleteApp = (token, id) => {
    return axios.post(`${BASE_URL}/fastgpt/apps/delete`, { id }, {
        headers: getAuthHeader(token)
    });
};

// Chat Interfaces
export const chatCompletion = async (data) => {
    const token = localStorage.getItem('token');
    if (data.stream) {
        // Return fetch promise for streaming - use stream endpoint
        return fetch(`${BASE_URL}/fastgpt/v1/chat/completions/stream`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...getAuthHeader(token)
            },
            body: JSON.stringify(data)
        });
    }
    // Standard axios for non-stream
    return axios.post(`${BASE_URL}/fastgpt/v1/chat/completions`, data, {
        headers: getAuthHeader(token)
    });
};

export const getHistories = (data) => {
    const token = localStorage.getItem('token');
    return axios.post(`${BASE_URL}/fastgpt/core/chat/history/getHistories`, data, {
        headers: getAuthHeader(token)
    });
};

export const getPaginationRecords = (data) => {
    const token = localStorage.getItem('token');
    return axios.post(`${BASE_URL}/fastgpt/core/chat/getPaginationRecords`, data, {
        headers: getAuthHeader(token)
    });
};

export const updateHistory = (data) => {
    const token = localStorage.getItem('token');
    return axios.post(`${BASE_URL}/fastgpt/core/chat/history/updateHistory`, data, {
        headers: getAuthHeader(token)
    });
};

export const delHistory = (appId, chatId) => {
    const token = localStorage.getItem('token');
    return axios.delete(`${BASE_URL}/fastgpt/core/chat/history/delHistory`, {
        params: { appId, chatId },
        headers: getAuthHeader(token)
    });
};

export const initOutLinkChat = (chatId, shareId, outLinkUid) => {
    const token = localStorage.getItem('token');
    return axios.get(`${BASE_URL}/fastgpt/core/chat/outLink/init`, {
        params: { chatId, shareId, outLinkUid },
        headers: getAuthHeader(token)
    });
};
