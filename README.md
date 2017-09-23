#Description
On Spotify there's no way to add all of an artist's albums, singles and features into a single playlist other than doing it manually. Spotify does have an "Add to playlist" option beside albums but this is only really useful for albums by that artist. In hip hop, artists often feature on other artists tracks and in that case, you don't want all tracks from the album, but just the song they feature on. Trying to create a discography for an artist like Drake or Future, who have hundreds of individual features, would be painful. Even worse would be someone like Gucci Mane, who has over *3000* tracks on Spotify.

This is where [Discogrify](https://discogrify.com/) steps in. Using the site, all you have to do is login via Spotify and search for an artist and you'll be presented with an artist's entire Spotify discography. All tracks from all albums, singles, features and compilations. You can optionally edit the selection of tracks to exclude individual tracks or entire albums. Once you're happy with the selection, simply publish the playlist and it will be in your library.

#Usage 
First things first, UI design is not my forte so I'm hoping for a pass on that one. Any suggestions are more than welcome. 
Secondly, ***If you notice anything wrong with the site, please leave a comment so I can fix the issue*** :)

Anyway, there are only a few steps to using the site:

1. Click the "Login via Spotify" button, which will redirect you to the Spotify authentication page.
2. Once authenticated, you'll be redirected back to Discogrify where you can perform a search. Results will appear in a list below the search box.
3. Once you select an artist from the list, the site will compile every track on Spotify containing that artist into a single list and display it in the results view.
4. In the results view, each row represents an album, be it by your searched artist or another artist. If you click on an album, it will expand display a list of tracks. You can select/deselect individual tracks or entire albums. This step is not mandatory
5. Selecting "Publish playlist" will create the playlist and add it to your account.  

# Future work
At the moment this is just a working prototype and if it's seeing good usage then I have a few improvements I'd like to make:

* **Popular artists**: Display a list of the most searched for artists
* **How-to/Guide/Landing page**: A brief intro about what the site does and how to use it.
* **Add multiple artists to one playlist**: The option to include multiple artists in one playlist. This would enable you to catalog, for example, all tracks by each DOOM moniker or all tracks by each individual member of Wu-Tang
* **Automatic filtering of duplicates and remixes**: I'm aware that there are duplicates in the results and the reason some albums/tracks are released differently in different regions and all version are returned by default on Spotify. Also the results don't consider a single that also exists on an album and will return them both by default.
Of course you can fix all of this manually in the results view but it would be nice to have it done automatically.
* **Alter existing user playlists**: I'd like to be able to add tracks to an existing playlist rather than one just created by the tool
* **Select/Deselect all**: A checkbox to select/deselect all results 
