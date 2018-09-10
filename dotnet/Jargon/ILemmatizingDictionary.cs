namespace Jargon
{
    public interface ILemmatizingDictionary
    {
        (string Canonical, bool Found) Lookup(string[] term);
    }
}
