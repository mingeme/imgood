<script lang='ts'>

  const { data } = $props();
  const { image } = $derived(data);
</script>

<div class='hero'>
  <div class='hero-body py-4'>
    <div class='container'>
      <div class='columns is-centered'>
        <div class='column is-10'>
          <div class='box p-0'>
            <img
              src={image.url}
              alt={image.name}
              style='object-fit: contain; max-height: 85dvh;'
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .hero-body {
    padding: 1rem 1.5rem;
  }

  .box {
    box-shadow: 0 2px 3px rgba(10, 10, 10, 0.1);
    transition: box-shadow 0.3s ease;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 200px;
  }

  .box:hover {
    box-shadow: 0 4px 6px rgba(10, 10, 10, 0.1);
  }

  img {
    max-width: 100%;
    height: auto;
  }
</style>
