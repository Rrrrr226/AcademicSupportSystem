import axios from 'axios';

export const getSubjectLink = (staffId) => {
  return axios.get(`/subject/get/links/${staffId}`);
};