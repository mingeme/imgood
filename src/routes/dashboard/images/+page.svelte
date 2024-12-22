<script lang='ts'>
  const { data } = $props();
  const { images } = $derived(data);
  function formatFileSize(bytes: number): string {
    if (bytes === 0)
      return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${Number.parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`;
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleString();
  }
</script>

<section class='section'>
  <div class='container'>
    <h1 class='title'>My Images</h1>

    <div class='table-container'>
      <table class='table is-fullwidth is-striped is-hoverable'>
        <thead>
          <tr>
            <th>Preview</th>
            <th>File Name</th>
            <th>Size</th>
            <th>Upload Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each images as image}
            <tr>
              <td>
                <img src={image.url} alt={image.name} />
              </td>
              <td>
                <a href={image.url}>{image.name}</a>
              </td>
              <td>{formatFileSize(image.file_size)}</td>
              <td>{formatDate(image.created_at)}</td>
              <td>
                <div class='buttons are-small'>
                  <button class='button is-info' onclick={() => window.open(`/api/images/${image.oss_key}`)}>
                    View Image
                  </button>
                  <button class='button is-primary' onclick={() => navigator.clipboard.writeText(`/api/images/${image.oss_key}`)}>
                    Copy URL
                  </button>
                  <button class='button is-danger'>
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</section>

<style>
  .table-container {
    overflow-x: auto;
  }

  .buttons {
    display: flex;
  }

</style>
