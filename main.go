// Read-only cgroup v2 resource explorer; never modifies controller files.
package main
import("encoding/json";"flag";"fmt";"os";"path/filepath";"strconv";"strings")
type Group struct{Path string `json:"path"`;MemoryCurrent *uint64 `json:"memory_current_bytes"`;MemoryMax string `json:"memory_max"`;CPUQuotaCores *float64 `json:"cpu_quota_cores"`;CPUUsageUsec *uint64 `json:"cpu_usage_usec"`;PidsCurrent *uint64 `json:"pids_current"`;PidsMax string `json:"pids_max"`;Errors []string `json:"read_errors,omitempty"`}
func quota(text string)(*float64,error){parts:=strings.Fields(text);if len(parts)!=2{return nil,fmt.Errorf("expected quota period")};period,e:=strconv.ParseUint(parts[1],10,64);if e!=nil||period==0{return nil,fmt.Errorf("invalid period")};if parts[0]=="max"{return nil,nil};n,e:=strconv.ParseUint(parts[0],10,64);if e!=nil{return nil,e};value:=float64(n)/float64(period);return &value,nil}
func readGroup(path string)Group{
 g:=Group{Path:path};read:=func(name string)string{b,e:=os.ReadFile(filepath.Join(path,name));if e!=nil{g.Errors=append(g.Errors,name+": "+e.Error());return ""};return strings.TrimSpace(string(b))}
 number:=func(name string)*uint64{text:=read(name);if text==""{return nil};v,e:=strconv.ParseUint(text,10,64);if e!=nil{g.Errors=append(g.Errors,name+": invalid number");return nil};return &v}
 g.MemoryCurrent=number("memory.current");g.MemoryMax=read("memory.max");g.PidsCurrent=number("pids.current");g.PidsMax=read("pids.max")
 if text:=read("cpu.max");text!=""{var e error;g.CPUQuotaCores,e=quota(text);if e!=nil{g.Errors=append(g.Errors,"cpu.max: "+e.Error())}}
 for _,line:=range strings.Split(read("cpu.stat"),"\n"){v:=strings.Fields(line);if len(v)==2&&v[0]=="usage_usec"{n,e:=strconv.ParseUint(v[1],10,64);if e==nil{g.CPUUsageUsec=&n}}}
 return g
}
func scan(root string,depth int)([]Group,error){
 root,e:=filepath.Abs(root);if e!=nil{return nil,e};info,e:=os.Stat(root);if e!=nil||!info.IsDir(){return nil,fmt.Errorf("root must be a directory")};groups:=[]Group{}
 e=filepath.WalkDir(root,func(path string,d os.DirEntry,err error)error{
  if err!=nil{if path==root{return err};return nil};if !d.IsDir(){return nil};rel,e:=filepath.Rel(root,path);if e!=nil{return e};level:=0;if rel!="."{level=len(strings.Split(rel,string(os.PathSeparator)))};if level>depth{return filepath.SkipDir}
  if len(groups)>=2000{return fmt.Errorf("group limit 2000 exceeded; choose a smaller root/depth")};groups=append(groups,readGroup(path));return nil
 });return groups,e
}
func main(){
 root:=flag.String("root","/sys/fs/cgroup","cgroup v2 directory or synthetic fixture root");depth:=flag.Int("depth",2,"maximum depth 0..8");jsonOutput:=flag.Bool("json",false,"emit structured JSON");flag.Parse()
 if *depth<0||*depth>8||flag.NArg()!=0{fmt.Fprintln(os.Stderr,"depth must be 0..8; use --help");os.Exit(2)}
 if _,e:=os.Stat(filepath.Join(*root,"cgroup.controllers"));e!=nil{fmt.Fprintln(os.Stderr,"root must expose cgroup.controllers (cgroup v2)");os.Exit(1)}
 groups,e:=scan(*root,*depth);if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 if *jsonOutput{encoder:=json.NewEncoder(os.Stdout);encoder.SetIndent("","  ");if e=encoder.Encode(groups);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)};return}
 number:=func(p *uint64)string{if p==nil{return "n/a"};return strconv.FormatUint(*p,10)}
 for _,g:=range groups{cpu:="unlimited/unknown";if g.CPUQuotaCores!=nil{cpu=fmt.Sprintf("%.2f cores",*g.CPUQuotaCores)};fmt.Printf("%s\n  memory %s / %s bytes · pids %s / %s · CPU %s · used %s µs\n",g.Path,number(g.MemoryCurrent),g.MemoryMax,number(g.PidsCurrent),g.PidsMax,cpu,number(g.CPUUsageUsec));if len(g.Errors)>0{fmt.Printf("  unavailable: %d fields (use --json for detail)\n",len(g.Errors))}}
}
