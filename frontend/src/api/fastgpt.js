import axios from 'axios';
import { BASE_URL } from './config';

const getAuthHeader = (token) => ({ Authorization: `Bearer ${token}` });

export const getAppList = (token, page = 1, pageSize = 10) => {
    return axios.post(`${BASE_URL}/fastgpt/apps/list`, {
        page,
        pageSize
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
