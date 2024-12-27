<script lang='ts'>
  import type { FileWithPreview } from '$lib/types';
  import { calculateFileHash } from '$lib/utils/file';

  const { user, files, onFilesSelected, onClear } = $props<{
    user: any;
    files: FileWithPreview[];
    onFilesSelected: (files: File[]) => void;
    onClear: () => void;
  }>();

  function removeFile(index: number) {
    URL.revokeObjectURL(files[index].preview);
    const newFiles = [...files];
    newFiles.splice(index, 1);
    onFilesSelected(newFiles.map(f => f.file));
  }

  async function upload() {
    for (const { file } of files) {
      const hash = await calculateFileHash(file);
      const res = await fetch('?/upload', {
        method: 'POST',
        body: JSON.stringify({
          name: file.name,
          size: file.size,
          hash,
        }),
      });
      const result = await res.json();
      if (result.type !== 'success') {
        console.error('Upload failed:', result.error);
        continue;
      }
      const { data } = result;
      await fetch(data.url, {
        method: 'PUT',
        body: file,
      });
    }
    window.location.reload();
  }
</script>

<div>
  <div class='file is-boxed file-cta image-upload-height'>
    <div
      class='container file-label is-responsive'
      class:cursor-default={!user}
    >
      <span class='image-upload-content'>
        {#if !user}
          <div class='has-text-centered'>
            You need to <a href='/signin'>sign in</a> before you can upload images
          </div>
        {:else if files.length > 0}
          {#each files as { preview, file }}
            <div class='preview-container'>
              <button
                type='button'
                class='delete-button'
                onclick={(e) => {
                  e.stopPropagation();
                  removeFile(files.findIndex((f: FileWithPreview) => f.preview === preview));
                }}
                aria-label='Delete image'
              >
                ×
              </button>
              <img
                src={preview}
                alt={file.name}
                class='preview-image'
              />
              <span class='file-name is-small' title={file.name}>
                {file.name}
              </span>
            </div>
          {/each}
        {:else}
          <div class='has-text-centered'>
            Drag &amp; drop files here ...
            <br />
            or
            <br />
            Copy &amp; paste screenshots here ...
          </div>
        {/if}
      </span>
    </div>
  </div>

  {#if user}
    <div class='is-flex'>
      <label class='file-label container select-files-label'>
        <input
          id='file-upload'
          class='file-input select-files-input'
          type='file'
          name='image'
          accept='.jpg, .jpeg, .png, .gif, .bmp'
          multiple
          onchange={(e: Event) => {
            const target = e.target as HTMLInputElement;
            onFilesSelected(Array.from(target.files ?? []));
          }}
        />
        <span class='select-files'>Select images...</span>
      </label>

      {#if files.length > 0}
        <div class='is-flex'>
          <button
            class='button is-danger select-button'
            type='button'
            onclick={onClear}
          >
            Clear
          </button>
          <button
            class='button is-primary select-button'
            onclick={upload}
          >
            Upload
          </button>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .preview-container {
    position: relative;
    display: inline-block;
    margin: 0.5rem;
    border-radius: 0.25rem;
    overflow: hidden;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .delete-button {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    width: 1.5rem;
    height: 1.5rem;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.5);
    color: white;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    transition: background-color 0.2s;
  }

  .delete-button:hover {
    background: rgba(0, 0, 0, 0.7);
  }

  .preview-image {
    width: 200px;
    height: 200px;
    object-fit: cover;
    display: block;
  }

  .image-upload-height {
    height: 100%;
    min-height: 20rem;
    border: 2px dashed var(--bulma-border);
    border-radius: 0.5rem;
    transition: border-color 0.2s;
  }

  .image-upload-height:hover {
    border-color: var(--bulma-primary);
  }

  .image-upload-content {
    margin-block: auto;
    overflow: auto;
    padding: 1rem;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    justify-content: center;
  }

  .select-files-label {
    border: 1px solid var(--bulma-border);
    border-radius: 0.25rem;
    overflow: hidden;
  }

  .select-files {
    display: flex;
    align-items: center;
    padding: 0.5rem 1rem;
  }

  .select-files-input {
    cursor: pointer;
  }

  .cursor-default {
    cursor: default;
  }

  .select-button {
    margin-left: 0.5rem;
  }
</style>
