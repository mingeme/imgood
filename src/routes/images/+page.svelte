<script lang='ts'>
  import { enhance } from '$app/forms';

  const { data } = $props();
  const { images } = $derived(data);

  const modalState = $state({ ossKey: '', isActive: false, isLoading: false });
  const copiedStates = $state<Record<string, boolean>>({});

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

  function openDeleteModal(ossKey: string) {
    modalState.ossKey = ossKey;
    modalState.isActive = true;
  }

  function closeDeleteModal() {
    modalState.ossKey = '';
    modalState.isActive = false;
    modalState.isLoading = false;
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
                <img src={image.preview_url} alt={image.name} />
              </td>
              <td>
                <a href={image.url}>{image.name}</a>
              </td>
              <td>{formatFileSize(image.file_size)}</td>
              <td>{formatDate(image.created_at)}</td>
              <td>
                <div class='buttons are-small'>
                  <a href={`/images/${image.id}`} class='button is-info'>
                    View Image
                  </a>
                  <button class='button is-primary' onclick={async () => {
                    await navigator.clipboard.writeText(image.url);
                    copiedStates[image.oss_key] = true;
                    setTimeout(() => copiedStates[image.oss_key] = false, 2000);
                  }}>
                    {#if copiedStates[image.oss_key]}
                      Copied!
                    {:else}
                      Copy URL
                    {/if}
                  </button>
                  <button class='button is-danger' onclick={() => openDeleteModal(image.oss_key)}>
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

<div class='modal' class:is-active={modalState.isActive}>
  <button
    class='modal-background'
    type='button'
    onclick={closeDeleteModal}
    aria-label='close modal'
    disabled={modalState.isLoading}
  ></button>
  <div class='modal-card'>
    <header class='modal-card-head'>
      <p class='modal-card-title'>Delete Image</p>
      <button class='delete' aria-label='close' onclick={closeDeleteModal}></button>
    </header>
    <section class='modal-card-body'>
      Are you sure you want to delete this image?
    </section>
    <footer class='modal-card-foot'>
      <form
        method='POST'
        action='?/delete'
        use:enhance={() => {
          modalState.isLoading = true;
          return async ({ result }) => {
            if (result.type === 'success') {
              window.location.reload();
            }
            else {
              modalState.isLoading = false;
            }
          };
        }}
      >
        <input type='hidden' name='ossKey' value={modalState.ossKey}>
        <button
          type='submit'
          class='button is-danger'
          disabled={modalState.isLoading}
          class:is-loading={modalState.isLoading}
        >
          Delete
        </button>
        <button type='button' class='button' onclick={closeDeleteModal}>Cancel</button>
      </form>
    </footer>
  </div>
</div>

<style>
  .table-container {
    overflow-x: auto;
  }

  .buttons {
    display: flex;
  }
</style>
