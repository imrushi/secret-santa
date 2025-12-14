<script>
  // Generate random snowflakes
  let flakes = $state(Array.from({ length: 50 }));
</script>

<div class="snow-container">
  {#each flakes as _, i}
    <div
      class="snowflake"
      style="
        left: {Math.random() * 100}vw;
        animation-duration: {Math.random() * 3 + 2}s;
        animation-delay: {Math.random() * 5}s;
        opacity: {Math.random()};
        font-size: {Math.random() * 10 + 10}px;
      "
    >
      ❄
    </div>
  {/each}
</div>

<style>
  .snow-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none; /* Let clicks pass through */
    z-index: 0;
    overflow: hidden;
  }

  .snowflake {
    position: absolute;
    top: -20px;
    color: white;
    user-select: none;
    animation-name: fall;
    animation-timing-function: linear;
    animation-iteration-count: infinite;
  }

  @keyframes fall {
    0% {
      transform: translateY(-20px) translateX(0px) rotate(0deg);
    }
    100% {
      transform: translateY(100vh) translateX(20px) rotate(360deg);
    }
  }
</style>
