import { Message } from '@arco-design/web-vue';

/**
 * ArcoMessage Hook
 */
export const useArcoMessage = () => {

  function successMessage(message: string) {
    Message.success(message);
  }

  function infoMessage(message: string) {
    Message.info(message);
  }

  function errorMessage(message: string) {
    Message.error(message);
  }

  function warnMessage(message: string) {
    Message.warning(message);
  }

  function loading(message: string) {
    Message.loading(message);
  }

  function clearMessage() {
    Message.clear();
  }

  return { successMessage, infoMessage, errorMessage, warnMessage, loading, clearMessage }
}