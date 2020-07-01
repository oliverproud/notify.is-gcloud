import axios from "axios";
import { catchAxiosError } from "./error";

axios.defaults.baseURL = "http://localhost:1323"

export const post = (url, data) => {
  return axios.post(url, data).catch(catchAxiosError)
}
