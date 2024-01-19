import axios from 'axios';

let host = process.env.REACT_APP_API_HOST_URL;

export const getLogTable = async (params) => {
    let url = `${host}/internal/logs?` + new URLSearchParams(params).toString();
    try{
        const res = await axios.get(url);
        return res.data;
    } catch (e) {
        console.log(e);
        return;
    }
}