namespace Jargon
{
    public interface Dictionary
    {
        (string Canonical, bool Found) Lookup(string term);
    }
}
