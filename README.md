chwallpaper
===========

A simple go program to change wallpaper upon workspace change in Linux.


Schema for JSON
--------------------

Within the data object is a list of Key -> [Wallpaper]\(Int -> [String]) mappings. Each key maps to a list, so that there can be a slideshow on each workspace.

Confirmed Desktop Environments
-------------------------------
Works on:
- GNOME
- Cinnamon

Known Problems
-------------------
Unity (at least on 14.04, possibly older) tells xprop that there is only one desktop, even with multiple workspaces, therefore the wallpaper will not change on workspace switch.
