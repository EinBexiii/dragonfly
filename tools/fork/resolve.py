import subprocess
def sections(s):
    out=[];i=0
    while True:
        a=s.find('<<<<<<<',i)
        if a<0: out.append(('t',s[i:]));break
        out.append(('t',s[i:a]))
        mid=s.index('=======',a);end=s.index('>>>>>>>',mid);eol=s.index('\n',end)+1
        out.append(('c',s[s.index('\n',a)+1:mid],s[mid+len('=======\n'):end]))
        i=eol
    return out
def write(p,pick):
    s=open(p).read()
    open(p,'w').write(''.join(x[1] if x[0]=='t' else pick(x[1],x[2]) for x in sections(s)))
def union(o,t):
    seen=set(l.strip() for l in o.split('\n'))
    ex='\n'.join(l for l in t.split('\n') if l.strip() and l.strip() not in seen)
    return o+(ex+'\n' if ex else '')
for p in subprocess.run(['git','diff','--name-only','--diff-filter=U'],capture_output=True,text=True).stdout.split():
    if p.endswith('passive.go'):
        write(p, lambda o,t: '''\t// Entities with a running fuse resend their state every quarter second so
\t// viewers see the fuse time progress.
\tif f := p.Fuse(); f >= 0 && f%(time.Second/4) == 0 {
\t\te.UpdateState()
''')
    elif p.endswith('movement.go'):
        write(p, lambda o,t: t if 'driven' in t else o)
        s=open(p).read().replace('\t\tif posChanged || driven {','\t\tif posChanged || rotChanged || driven {',1)
        open(p,'w').write(s)
    elif p.endswith('player.go'):
        # keep whole func blocks from the incoming side, minus the two the UI
        # facade already owns. Dropping single lines leaves broken bodies.
        def keepfuncs(o, t):
            out, block, drop = [], [], False
            for line in t.split('\n'):
                if line.startswith('func ') or (line.startswith('//') and not block):
                    if block and not drop:
                        out.append('\n'.join(block))
                    block, drop = [], False
                if 'HideEntity' in line or 'ShowEntity' in line:
                    drop = True
                block.append(line)
            if block and not drop:
                out.append('\n'.join(block))
            return '\n'.join(out).strip('\n') + '\n'
        write(p, keepfuncs)
    else:
        write(p, union)
