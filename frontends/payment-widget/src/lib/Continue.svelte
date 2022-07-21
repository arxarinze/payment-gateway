<script>
  import { fly } from "svelte/transition";
  export let selectedPage, pages, data, text;
  function changePage() {
    let amount = data.amount;
    if (
      parseFloat(amount.fiat) != NaN &&
      parseFloat(amount.coin) != NaN &&
      amount.coin > 0 &&
      amount.fiat > 0
    ) {
      if (selectedPage < pages - 1) {
        selectedPage = selectedPage + 1;
      }
      if (text == "Start Over") {
        selectedPage = 0;
      }
    }
  }
</script>

<div
  class="btn-container"
  out:fly={{ x: 200, duration: 500 }}
  in:fly={{ x: -200, duration: 500 }}
>
  <button
    class="btn"
    style="
    margin: 2rem;"
    on:click={() => {
      changePage();
    }}>{text}</button
  >
</div>

<style>
  .btn {
    color: white;
    background-color: transparent;
    outline: 0;
    font-family: "Ubuntu", sans-serif;
    border: 0;
    font-size: larger;
  }
  .btn-container {
    background: black;
    position: fixed;
    width: 100vw;
    left: 0;
    bottom: 0;
    height: auto;
    right: 0;
  }
  .btn-container:hover {
    background-color: #fe3e00;
    transition: all;
  }
</style>
