import axios from "axios";
import config from "../config"; // Import URL từ file config
const API_URL = config.apiBaseUrl;

export const getDevices = () => axios.get(`${API_URL}/devices`);
export const getDeviceFunctions = (id) =>
  axios.get(`${API_URL}/device/${id}/functions`);
export const sendDeviceCommand = (id, data) =>
  axios.post(`${API_URL}/device/${id}/command`, data);
export const sendTurnOn = () => axios.get(`${API_URL}/device/on`);
export const sendTurnOff = () => axios.get(`${API_URL}/device/off`);

export const mapSwitch = (payload) =>
  axios.post(`${API_URL}/map-switch`, payload);
export const getMappedSwitches = () => axios.get(`${API_URL}/mapped-switches`);
// Update mapped switch
export const updateMappedSwitch = (id, payload) =>
  axios.put(`${API_URL}/map-switch/${id}`, payload);

// Delete mapped switch
export const deleteMappedSwitch = (id) =>
  axios.delete(`${API_URL}/map-switch/${id}`);

export const resetSwitchState = () =>
  axios.post(`${API_URL}/reset-switch-state`);

export const fetchSettings = () => axios.get(`${API_URL}/settings`);

export const updateSettings = (data) => axios.post(`${API_URL}/settings`, data); // gửi mảng [{ key, value }]
