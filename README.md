# Matrix.Nucleus

The server-side software for Matrix.Marketplace.
Which is a software distribution marketplace that made from developer & for developer.

> [!NOTE]
> Matrix is now a part of the HyperNet project. To run it yourself, you should run the service gateway Nexus and its dependencies.
> Learn more from Nexus project here: <https://github.com/Solsynth/HyperNet.Nexus>

## Concepts

One of the core concepts of Matrix is the **Product**.
We call a single application, or a game or something else a product.

When you want to release something like executable or something else.
You will need to create a **Release**.

The release is belongs to a product, which a product has multiple releases.

The release has a **Version**. Which follow the [Semantic Versioning](https://semver.org)
for better understanding and version comparison.

The release has also a channel, which different channel subscribers will receive different releases.
Like the beta, canary and the stable channel. You can add anything you want.

For user to download the release, there has an **Assets** field on the release.
The assets is a Dictionary of the files to download for different platform, the key is the platform id.
And the value is list of the compressed assets files or links.

> The file serving is powered by our [Paperclip](https://github.com/Solsynth/HyperNet.Paperclip) project. We're really proud of it.

### Updates

About the release, we got two types of the release, full release and increment one.

For new users, it will need to download the latest full release and all the increment releases after it,
and our app will merge all the release into a working one.

For the old users that already have the previous release, it will download the latest increment one
and apply it to the local one.
