# WebRT img2webp

Convertor de imagini în WebP, în linie de comandă. Scanează un director și convertește toate imaginile suportate (`.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.tif`, `.tiff`) în `.webp`, păstrând numele fișierului.

## Compilare

```bash
cd img2webp
go build -o img2webp .
```

## Utilizare

Rulează din directorul care conține imaginile (procesează directorul curent `.`):

```bash
./img2webp
```

## Opțiuni

| Flag | Implicit | Descriere |
|---|---|---|
| `-q` | `85` | Calitate WebP (0-100) |
| `-r` | `false` | Procesează și subdirectoarele, recursiv |
| `-delete` | `false` | Șterge fișierul original după conversia reușită |
| `-overwrite` | `false` | Suprascrie fișierele `.webp` existente (implicit sunt sărite) |

## Exemple

```bash
# conversie simplă, calitate implicită
./img2webp

# recursiv, calitate mai mare, șterge originalele
./img2webp -r -q 90 -delete

# reconvertește tot, chiar dacă .webp există deja
./img2webp -overwrite
```

## Note

- Fișierele `.webp` deja existente ca sursă sunt sărite automat (nu se convertesc la ele însele).
- Scrierea e sigură: se scrie mai întâi într-un fișier temporar, apoi se redenumește — nu rămân fișiere `.webp` corupte dacă encodarea eșuează la jumătate.
- La final se afișează un rezumat: câte fișiere au fost convertite, sărite sau eșuate.

## Instalare globală (`webrt-img2webp` / `webrt:img2webp`, în orice terminal)

Binarul e copiat în `/opt/webrt-tools/bin`, iar în `/usr/local/bin` (care e în `PATH` pe orice distribuție Linux, indiferent de shell) se creează symlink-uri sub ambele nume, cu liniuță și cu două puncte:

```bash
cd img2webp
go build -o img2webp .

sudo mkdir -p /opt/webrt-tools/bin
sudo install -m 755 img2webp /opt/webrt-tools/bin/img2webp

sudo ln -sf /opt/webrt-tools/bin/img2webp /usr/local/bin/webrt-img2webp
sudo ln -sf /opt/webrt-tools/bin/img2webp /usr/local/bin/webrt:img2webp
```

După asta, din orice director, în orice terminal:

```bash
webrt-img2webp -q 90 -r
# sau
webrt:img2webp -q 90 -r
```

Dezinstalare:

```bash
sudo rm -f /usr/local/bin/webrt-img2webp /usr/local/bin/webrt:img2webp /opt/webrt-tools/bin/img2webp
```
